package com.blackmidori.familyexpenses.api.service;

import com.blackmidori.familyexpenses.api.model.Expense;
import com.blackmidori.familyexpenses.api.repository.ExpenseRepository;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;

import java.util.Optional;

@Service
public class ExpenseService {

    private final ExpenseRepository expenseRepository;

    public ExpenseService(ExpenseRepository expenseRepository) {
        this.expenseRepository = expenseRepository;
    }

    public Expense insert(Expense expense) {
        return this.expenseRepository.insert(expense);
    }

    public Page<Expense> findAll(Pageable pageable) {
        return expenseRepository.findAll(pageable);
    }
    public Optional<Expense> findById(String expenseId) {
        return expenseRepository.findById(expenseId);
    }

    public Expense update(Expense expense) {
        Assert.isTrue(this.expenseRepository.existsById(expense.getId()),"Expense not found for updating");
        return this.expenseRepository.save(expense);
    }
    public boolean existsById(String expenseId) {
        return expenseRepository.existsById(expenseId);
    }

    public void deleteById(String expenseId) {
        expenseRepository.deleteById(expenseId);
    }
}