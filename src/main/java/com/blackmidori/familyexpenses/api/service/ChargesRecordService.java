package com.blackmidori.familyexpenses.api.service;

import com.blackmidori.familyexpenses.api.application.exception.EntityNotFound;
import com.blackmidori.familyexpenses.api.dto.ChargesRecordDto;
import com.blackmidori.familyexpenses.api.factory.ChargesRecordFactory;
import com.blackmidori.familyexpenses.api.model.ChargesRecord;
import com.blackmidori.familyexpenses.api.repository.ChargesRecordRepository;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;

import java.util.Optional;

@Service
public class ChargesRecordService {

    private final ChargesRecordRepository chargesRecordRepository;
    private final ChargesRecordFactory chargesRecordFactory;

    public ChargesRecordService(ChargesRecordRepository chargesRecordRepository, ChargesRecordFactory chargesRecordFactory) {
        this.chargesRecordRepository = chargesRecordRepository;
        this.chargesRecordFactory = chargesRecordFactory;
    }


    public ChargesRecord convert(ChargesRecordDto chargesRecordDto) throws EntityNotFound {
        return chargesRecordFactory.createFromDto(chargesRecordDto);
    }

    public ChargesRecord insert(ChargesRecord chargesRecord) {
        return this.chargesRecordRepository.insert(chargesRecord);
    }

    public Page<ChargesRecord> findAll(Pageable pageable) {
        return chargesRecordRepository.findAll(pageable);
    }
    public Optional<ChargesRecord> findById(String chargesRecordId) {
        return chargesRecordRepository.findById(chargesRecordId);
    }

    public ChargesRecord update(ChargesRecord chargesRecord) {
        Assert.isTrue(this.chargesRecordRepository.existsById(chargesRecord.getId()),"ChargesRecord not found for updating");
        return this.chargesRecordRepository.save(chargesRecord);
    }
    public boolean existsById(String chargesRecordId) {
        return chargesRecordRepository.existsById(chargesRecordId);
    }

    public void deleteById(String chargesRecordId) {
        chargesRecordRepository.deleteById(chargesRecordId);
    }
}