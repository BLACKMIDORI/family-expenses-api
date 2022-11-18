package com.blackmidori.familyexpenses.api.repository;

import com.blackmidori.familyexpenses.api.model.ChargesModel;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface ChargesModelRepository extends MongoRepository<ChargesModel, String> {
}
