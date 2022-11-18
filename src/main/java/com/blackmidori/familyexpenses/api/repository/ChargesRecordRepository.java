package com.blackmidori.familyexpenses.api.repository;

import com.blackmidori.familyexpenses.api.model.ChargesRecord;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface ChargesRecordRepository extends MongoRepository<ChargesRecord, String> {
}
