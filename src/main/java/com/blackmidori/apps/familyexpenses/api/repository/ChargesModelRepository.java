package com.blackmidori.apps.familyexpenses.api.repository;

import com.blackmidori.apps.familyexpenses.api.model.ChargesModel;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface ChargesModelRepository extends MongoRepository<ChargesModel, String> {
}
