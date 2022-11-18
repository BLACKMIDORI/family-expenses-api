package com.blackmidori.familyexpenses.api.controller;

import com.blackmidori.familyexpenses.api.dto.PayerDto;
import com.blackmidori.familyexpenses.api.factory.PayerFactory;
import com.blackmidori.familyexpenses.api.model.Payer;
import com.blackmidori.familyexpenses.api.model.Workspace;
import com.blackmidori.familyexpenses.api.service.PayerService;
import com.blackmidori.familyexpenses.api.service.WorkspaceService;
import com.blackmidori.familyexpenses.api.util.UriUtils;
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
@RequestMapping("/payer")
public class PayerController {
    private static final Logger logger = Logger.getLogger(PayerController.class);

    private final PayerService payerService;
    private final WorkspaceService workspaceService;

    public PayerController(PayerService payerService, WorkspaceService workspaceService) {
        this.payerService = payerService;
        this.workspaceService = workspaceService;
    }

    @PostMapping
    @ResponseBody
    public ResponseEntity<Object> create(@RequestBody @Valid PayerDto payerDto){
        Optional<Workspace> workspaceOptional =workspaceService.findById(payerDto.getWorkspaceId());
        if(workspaceOptional.isEmpty()){
            return ResponseEntity.unprocessableEntity().body("incorrect workspaceId");
        }
        final var payer = new PayerFactory().createFromDto(payerDto, workspaceOptional.get());
        payerService.insert(payer);
        return ResponseEntity.created(UriUtils.getCreatedUrl(payer.getId())).body(payer);
    }


    @GetMapping
    public ResponseEntity<Page<Payer>> getAll(@PageableDefault(size = 10) Pageable pageable) {
        return ResponseEntity.ok(payerService.findAll(pageable));
    }


    @PutMapping("/{payerId}")
    public ResponseEntity<Object> update(@PathVariable String payerId,@RequestBody @Valid PayerDto payerDto) {
        final Optional<Payer> payerOptional = payerService.findById(payerId);
        if(payerOptional.isEmpty()){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Payer not found");
        }
        final Payer storedPayer = payerOptional.get();
        Optional<Workspace> workspace =workspaceService.findById(payerDto.getWorkspaceId());
        if(workspace.isEmpty()){
            return ResponseEntity.unprocessableEntity().body("incorrect workspaceId");
        }
        final var updatedPayer = new PayerFactory().createFromDto(payerDto, workspace.get());
        updatedPayer.setId(storedPayer.getId());
        updatedPayer.setCreationDateTime(storedPayer.getCreationDateTime());
        return ResponseEntity.ok(payerService.update(updatedPayer));
    }

    @DeleteMapping("/{payerId}")
    public ResponseEntity<Object> delete(@PathVariable String payerId) {
        if(!payerService.existsById(payerId)){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Payer not found");
        }
        payerService.deleteById(payerId);
        return ResponseEntity.noContent().build();
    }
}