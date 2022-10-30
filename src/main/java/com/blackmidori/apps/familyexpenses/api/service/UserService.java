package com.blackmidori.apps.familyexpenses.api.service;

import com.blackmidori.apps.familyexpenses.api.model.User;
import com.blackmidori.apps.familyexpenses.api.repository.UserRepository;
import org.springframework.stereotype.Service;

@Service
public class UserService {

    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public User findDevelopmentUser() {
        return userRepository.findDevelopmentUser();
    }
}