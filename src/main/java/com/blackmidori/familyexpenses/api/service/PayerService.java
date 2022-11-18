package com.blackmidori.familyexpenses.api.service;

import com.blackmidori.familyexpenses.api.model.Payer;
import com.blackmidori.familyexpenses.api.repository.PayerRepository;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;

import java.util.Optional;

@Service
public class PayerService {

    private final PayerRepository payerRepository;

    public PayerService(PayerRepository payerRepository) {
        this.payerRepository = payerRepository;
    }

    public Payer insert(Payer payer) {
        return this.payerRepository.insert(payer);
    }

    public Page<Payer> findAll(Pageable pageable) {
        return payerRepository.findAll(pageable);
    }
    public Optional<Payer> findById(String payerId) {
        return payerRepository.findById(payerId);
    }

    public Payer update(Payer payer) {
        Assert.isTrue(this.payerRepository.existsById(payer.getId()),"Payer not found for updating");
        return this.payerRepository.save(payer);
    }
    public boolean existsById(String payerId) {
        return payerRepository.existsById(payerId);
    }

    public void deleteById(String payerId) {
        payerRepository.deleteById(payerId);
    }
}