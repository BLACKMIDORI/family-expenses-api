package com.blackmidori.apps.familyexpenses.api.controller;

import com.blackmidori.apps.familyexpenses.api.model.User;
import com.blackmidori.apps.familyexpenses.api.model.Workspace;
import com.blackmidori.apps.familyexpenses.api.service.UserService;
import com.blackmidori.apps.familyexpenses.api.service.WorkspaceService;
import com.blackmidori.apps.familyexpenses.api.util.UriUtils;
import org.apache.log4j.Logger;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/workspace")
public class WorkspaceController {
    private static final Logger logger = Logger.getLogger(WorkspaceController.class);

    private final WorkspaceService workspaceService;
    private final UserService userService;

    public WorkspaceController(WorkspaceService workspaceService, UserService userService) {
        this.workspaceService = workspaceService;
        this.userService = userService;
    }


    @PostMapping
    @ResponseBody
    public ResponseEntity<Workspace> create() {
        final User user = userService.findDevelopmentUser();
        final Workspace workspace = new Workspace(user);
        workspaceService.insert(workspace);
        return ResponseEntity.created(UriUtils.getCreatedUrl(workspace.getId())).body(workspace);
    }
    @GetMapping
    public ResponseEntity<Page<Workspace>> getAll(@PageableDefault(size = 10) Pageable pageable) {
        final Page<Workspace> list = workspaceService.findAll(pageable);
        return ResponseEntity.ok(list);
    }

//    @PutMapping("/{workspace}")
//    public ResponseEntity<Object> update(@PathVariable String workspaceId) {
//        final Optional<Workspace> optionalWorkspace = workspaceService.findById(workspace);
//        if(optionalWorkspace.isEmpty()){
//            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Workspace not found");
//        }
//        final Workspace storedWorkspace = optionalWorkspace.get();
//        final var updatedPayer = new WorkspaceFactory().create();
//        updatedPayer.setId(storedPayer.getId());
//        updatedPayer.setCreationDateTime(storedPayer.getCreationDateTime());
//        return ResponseEntity.ok(payerService.update(updatedPayer));
//    }


    @DeleteMapping("{workspaceId}")
    public ResponseEntity<Object> delete(@PathVariable String workspaceId) {
        if(!workspaceService.existsById(workspaceId)){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Workspace not found");
        }
        workspaceService.deleteById(workspaceId);
        return ResponseEntity.noContent().build();
    }

}