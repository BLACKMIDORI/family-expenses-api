package com.blackmidori.familyexpenses.api.repository;

import com.blackmidori.familyexpenses.api.model.Expense;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface ExpenseRepository extends MongoRepository<Expense, String> {
}
