package com.blackmidori.familyexpenses.api.repository;

import com.blackmidori.familyexpenses.api.model.User;
import org.springframework.data.mongodb.repository.MongoRepository;

public interface UserRepository extends MongoRepository<User, String>{
    default User findDevelopmentUser(){
        return User.developmentUser;
    }
}
