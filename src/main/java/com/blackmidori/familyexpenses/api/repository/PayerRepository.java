package com.blackmidori.familyexpenses.api.repository;

import com.blackmidori.familyexpenses.api.model.Payer;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface PayerRepository extends MongoRepository<Payer, String> {
}
