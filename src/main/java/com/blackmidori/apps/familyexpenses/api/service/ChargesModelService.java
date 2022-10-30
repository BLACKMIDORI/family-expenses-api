package com.blackmidori.apps.familyexpenses.api.service;

import com.blackmidori.apps.familyexpenses.api.application.exception.EntityNotFound;
import com.blackmidori.apps.familyexpenses.api.dto.ChargesModelDto;
import com.blackmidori.apps.familyexpenses.api.factory.ChargesModelFactory;
import com.blackmidori.apps.familyexpenses.api.model.ChargesModel;
import com.blackmidori.apps.familyexpenses.api.repository.ChargesModelRepository;
import com.blackmidori.apps.familyexpenses.api.repository.WorkspaceRepository;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;

import java.util.Optional;

@Service
public class ChargesModelService {

    private final ChargesModelRepository chargesModelRepository;
    private final ChargesModelFactory chargesModelFactory;

    public ChargesModelService(ChargesModelRepository chargesModelRepository, ChargesModelFactory chargesModelFactory) {
        this.chargesModelRepository = chargesModelRepository;
        this.chargesModelFactory = chargesModelFactory;
    }

    public ChargesModel convert(ChargesModelDto chargesModelDto) throws EntityNotFound {
        return chargesModelFactory.createFromDto(chargesModelDto);
    }

    public ChargesModel insert(ChargesModel chargesModel) {
        return this.chargesModelRepository.insert(chargesModel);
    }

    public Page<ChargesModel> findAll(Pageable pageable) {
        return chargesModelRepository.findAll(pageable);
    }
    public Optional<ChargesModel> findById(String chargesModelId) {
        return chargesModelRepository.findById(chargesModelId);
    }

    public ChargesModel update(ChargesModel chargesModel) {
        Assert.isTrue(this.chargesModelRepository.existsById(chargesModel.getId()),"ChargesModel not found for updating");
        return this.chargesModelRepository.save(chargesModel);
    }
    public boolean existsById(String chargesModelId) {
        return chargesModelRepository.existsById(chargesModelId);
    }

    public void deleteById(String chargesModelId) {
        chargesModelRepository.deleteById(chargesModelId);
    }
}