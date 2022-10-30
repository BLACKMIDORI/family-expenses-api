package com.blackmidori.apps.familyexpenses.api.repository;

import com.blackmidori.apps.familyexpenses.api.model.Payer;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface PayerRepository extends MongoRepository<Payer, String> {
}
