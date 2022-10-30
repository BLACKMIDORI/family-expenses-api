package com.blackmidori.apps.familyexpenses.api.repository;

import com.blackmidori.apps.familyexpenses.api.model.User;
import com.blackmidori.apps.familyexpenses.api.model.Workspace;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.mongodb.repository.Query;

import java.util.List;

public interface WorkspaceRepository extends MongoRepository<Workspace, String> {
    @Query("{'owner' : ?0}")
    List<Workspace> findAllByUser(User user);
}
