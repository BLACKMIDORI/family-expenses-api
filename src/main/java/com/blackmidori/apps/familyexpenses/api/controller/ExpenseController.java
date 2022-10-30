package com.blackmidori.apps.familyexpenses.api.controller;

import com.blackmidori.apps.familyexpenses.api.dto.ExpenseDto;
import com.blackmidori.apps.familyexpenses.api.dto.PayerDto;
import com.blackmidori.apps.familyexpenses.api.factory.ExpenseFactory;
import com.blackmidori.apps.familyexpenses.api.factory.PayerFactory;
import com.blackmidori.apps.familyexpenses.api.model.Expense;
import com.blackmidori.apps.familyexpenses.api.model.Payer;
import com.blackmidori.apps.familyexpenses.api.model.Workspace;
import com.blackmidori.apps.familyexpenses.api.service.ExpenseService;
import com.blackmidori.apps.familyexpenses.api.service.PayerService;
import com.blackmidori.apps.familyexpenses.api.service.WorkspaceService;
import com.blackmidori.apps.familyexpenses.api.util.UriUtils;
import org.apache.log4j.Logger;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.Optional;

@RestController
@RequestMapping("/expense")
public class ExpenseController {
    private static final Logger logger = Logger.getLogger(ExpenseController.class);

    private final ExpenseService expenseService;
    private final WorkspaceService workspaceService;

    public ExpenseController(ExpenseService expenseService, WorkspaceService workspaceService) {
        this.expenseService = expenseService;
        this.workspaceService = workspaceService;
    }

    @PostMapping
    @ResponseBody
    public ResponseEntity<Object> create(@RequestBody @Valid ExpenseDto expenseDto){
        Optional<Workspace> workspace =workspaceService.findById(expenseDto.getWorkspaceId());
        if(workspace.isEmpty()){
            return ResponseEntity.unprocessableEntity().body("incorrect workspaceId");
        }
        final var expense = new ExpenseFactory().createFromDto(expenseDto, workspace.get());
        expenseService.insert(expense);
        return ResponseEntity.created(UriUtils.getCreatedUrl(expense.getId())).body(expense);
    }


    @GetMapping
    public ResponseEntity<Page<Expense>> getAll(@PageableDefault(size = 10) Pageable pageable) {
        return ResponseEntity.ok(expenseService.findAll(pageable));
    }


    @PutMapping("/{expenseId}")
    public ResponseEntity<Object> update(@PathVariable String expenseId,@RequestBody @Valid ExpenseDto expenseDto) {
        final Optional<Expense> expenseOptional = expenseService.findById(expenseId);
        if(expenseOptional.isEmpty()){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Expense not found");
        }
        final Expense storedExpense = expenseOptional.get();
        Optional<Workspace> workspaceOptional =workspaceService.findById(expenseDto.getWorkspaceId());
        if(workspaceOptional.isEmpty()){
            return ResponseEntity.unprocessableEntity().body("incorrect workspaceId");
        }
        final var updatedExpense = new ExpenseFactory().createFromDto(expenseDto, workspaceOptional.get());
        updatedExpense.setId(storedExpense.getId());
        updatedExpense.setCreationDateTime(storedExpense.getCreationDateTime());
        return ResponseEntity.ok(expenseService.update(updatedExpense));
    }

    @DeleteMapping("/{expenseId}")
    public ResponseEntity<Object> delete(@PathVariable String expenseId) {
        if(!expenseService.existsById(expenseId)){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Expense not found");
        }
        expenseService.deleteById(expenseId);
        return ResponseEntity.noContent().build();
    }
}