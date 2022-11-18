package com.blackmidori.familyexpenses.api.controller;

import com.blackmidori.familyexpenses.api.application.exception.EntityNotFound;
import com.blackmidori.familyexpenses.api.dto.ChargesRecordDto;
import com.blackmidori.familyexpenses.api.model.ChargesRecord;
import com.blackmidori.familyexpenses.api.service.ChargesRecordService;
import com.blackmidori.familyexpenses.api.util.UriUtils;
import org.apache.log4j.Logger;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.time.OffsetDateTime;
import java.util.Optional;

@RestController
@RequestMapping("/charges-record")
public class ChargesRecordController {
    private static final Logger logger = Logger.getLogger(ChargesRecordController.class);

    private final ChargesRecordService chargesRecordService;

    public ChargesRecordController(ChargesRecordService chargesRecordService) {
        this.chargesRecordService = chargesRecordService;
    }

    @PostMapping
    @ResponseBody
    public ResponseEntity<Object> create(@RequestBody @Valid ChargesRecordDto chargesRecordDto){
        final ChargesRecord chargesRecord;
        try {
            chargesRecord = chargesRecordService.convert(chargesRecordDto);
            chargesRecord.setOpeningDateTime(OffsetDateTime.now());
        } catch (EntityNotFound e) {
            return ResponseEntity.unprocessableEntity().body(e.getMessage());
        }
        chargesRecordService.insert(chargesRecord);
        return ResponseEntity.created(UriUtils.getCreatedUrl(chargesRecord.getId())).body(chargesRecord);
    }


    @GetMapping
    public ResponseEntity<Page<ChargesRecord>> getAll(@PageableDefault(size = 10) Pageable pageable) {
        return ResponseEntity.ok(chargesRecordService.findAll(pageable));
    }


    @PutMapping("/{chargesRecordId}")
    public ResponseEntity<Object> update(@PathVariable String chargesRecordId,@RequestBody @Valid ChargesRecordDto chargesRecordDto) {
        final Optional<ChargesRecord> chargesRecordOptional = chargesRecordService.findById(chargesRecordId);
        if(chargesRecordOptional.isEmpty()){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("ChargesModel not found");
        }
        final ChargesRecord storedChargesRecord = chargesRecordOptional.get();

        final ChargesRecord updatedChargesRecord;
        try {
            updatedChargesRecord = chargesRecordService.convert(chargesRecordDto);
        } catch (EntityNotFound e) {
            return ResponseEntity.unprocessableEntity().body(e.getMessage());
        }
        updatedChargesRecord.setId(storedChargesRecord.getId());
        updatedChargesRecord.setCreationDateTime(storedChargesRecord.getCreationDateTime());
        updatedChargesRecord.setOpeningDateTime(storedChargesRecord.getOpeningDateTime());
        return ResponseEntity.ok(chargesRecordService.update(updatedChargesRecord));
    }

    @DeleteMapping("/{chargesRecordId}")
    public ResponseEntity<Object> delete(@PathVariable String chargesRecordId) {
        if(!chargesRecordService.existsById(chargesRecordId)){
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body("ChargesRecord not found");
        }
        chargesRecordService.deleteById(chargesRecordId);
        return ResponseEntity.noContent().build();
    }
}