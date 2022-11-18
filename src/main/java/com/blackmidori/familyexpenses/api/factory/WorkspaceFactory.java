package com.blackmidori.familyexpenses.api.factory;

import com.blackmidori.familyexpenses.api.model.User;
import com.blackmidori.familyexpenses.api.model.Workspace;
import lombok.NonNull;

public class WorkspaceFactory {

    public Workspace create(@NonNull User user) {
        return new Workspace(user);
    }
}
