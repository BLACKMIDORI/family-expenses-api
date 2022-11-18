package com.blackmidori.familyexpenses.api.service;

import com.blackmidori.familyexpenses.api.model.User;
import com.blackmidori.familyexpenses.api.model.Workspace;
import com.blackmidori.familyexpenses.api.repository.UserRepository;
import com.blackmidori.familyexpenses.api.repository.WorkspaceRepository;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class WorkspaceService {

    private final WorkspaceRepository workspaceRepository;
    private final UserRepository userRepository;

    public WorkspaceService(WorkspaceRepository workspaceRepository, UserRepository userRepository) {
        this.workspaceRepository = workspaceRepository;
        this.userRepository = userRepository;
    }

    public Workspace insert(Workspace workspace) {
        return workspaceRepository.insert(workspace);
    }

    public Page<Workspace> findAll(Pageable pageable) {
        return workspaceRepository.findAll(pageable);
    }
    public List<Workspace> find(User user) {
        return workspaceRepository.findAllByUser(user);
    }

    public Optional<Workspace> findById(String workspaceId) {
        return workspaceRepository.findById(workspaceId);
    }
    public boolean existsById(String workspaceId) {
        return workspaceRepository.existsById(workspaceId);
    }public void deleteById(String workspaceId) {
        workspaceRepository.deleteById(workspaceId);
    }
}