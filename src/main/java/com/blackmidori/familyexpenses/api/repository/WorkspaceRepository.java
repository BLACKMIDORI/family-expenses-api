package com.blackmidori.familyexpenses.api.repository;

import com.blackmidori.familyexpenses.api.model.User;
import com.blackmidori.familyexpenses.api.model.Workspace;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.mongodb.repository.Query;

import java.util.List;

public interface WorkspaceRepository extends MongoRepository<Workspace, String> {
    @Query("{'owner' : ?0}")
    List<Workspace> findAllByUser(User user);
}
