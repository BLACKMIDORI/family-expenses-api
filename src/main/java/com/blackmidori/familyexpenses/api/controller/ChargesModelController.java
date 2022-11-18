package com.blackmidori.familyexpenses.api.controller;

import com.blackmidori.familyexpenses.api.application.exception.EntityNotFound;
import com.blackmidori.familyexpenses.api.dto.ChargesModelDto;
import com.blackmidori.familyexpenses.api.model.ChargesModel;
import com.blackmidori.familyexpenses.api.service.ChargesModelService;
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
@RequestMapping("/charges-model")
public class ChargesModelController {
    private static final Logger logger = Logger.getLogger(ChargesModelController.class);

    private final ChargesModelService chargesModelService;

    public ChargesModelController(ChargesModelService chargesModelService) {
        this.chargesModelService = chargesModelService;
    }

    @PostMapping
    @ResponseBody
    public ResponseEntity<Object> create(@RequestBody @Valid ChargesModelDto chargesModelDto){
        final ChargesModel chargesModel;
        try {
            chargesModel = chargesModelService.convert(chargesModelDto);
            chargesModel.setUpdatingDateTime(chargesModel.getCreationDateTime());
        } catch (EntityNotFound e) {
            return ResponseEntity.unprocessableEntity().body(e.getMessage());
        }
        chargesModelService.insert(chargesModel);
        return ResponseEntity.created(UriUtils.getCreatedUrl(chargesModel.getId())).body(chargesModel);
    }


    @GetMapping
    public ResponseEntity<Page<ChargesModel>> getAll(@PageableDefault(size = 10) Pageable pageable) {
        return ResponseEntity.ok(chargesModelService.findAll(pageable));
    }


    @PutMapping("/{chargesModelId}")
    public ResponseEntity<Object> update(@PathVariable String chargesModelId,@RequestBody @Valid ChargesModelDto chargesModelDto) {
        final Optional<ChargesModel> chargesModelOptional = chargesModelService.findById(chargesModelId);
        if(chargesModelOptional.isEmpty()){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("ChargesModel not found");
        }
        final ChargesModel storedChargesModel = chargesModelOptional.get();

        final ChargesModel updatedChargesModel;
        try {
            updatedChargesModel = chargesModelService.convert(chargesModelDto);
        } catch (EntityNotFound e) {
            return ResponseEntity.unprocessableEntity().body(e.getMessage());
        }
        updatedChargesModel.setId(storedChargesModel.getId());
        updatedChargesModel.setCreationDateTime(storedChargesModel.getCreationDateTime());
        return ResponseEntity.ok(chargesModelService.update(updatedChargesModel));
    }

    @DeleteMapping("/{chargesModelId}")
    public ResponseEntity<Object> delete(@PathVariable String chargesModelId) {
        if(!chargesModelService.existsById(chargesModelId)){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("ChargesModel not found");
        }
        chargesModelService.deleteById(chargesModelId);
        return ResponseEntity.noContent().build();
    }
}