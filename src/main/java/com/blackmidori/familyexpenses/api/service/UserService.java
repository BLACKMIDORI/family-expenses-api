package com.blackmidori.familyexpenses.api.service;

import com.blackmidori.familyexpenses.api.model.User;
import com.blackmidori.familyexpenses.api.repository.UserRepository;
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