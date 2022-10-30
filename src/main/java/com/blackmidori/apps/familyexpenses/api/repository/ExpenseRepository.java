package com.blackmidori.apps.familyexpenses.api.repository;

import com.blackmidori.apps.familyexpenses.api.model.Expense;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface ExpenseRepository extends MongoRepository<Expense, String> {
}
