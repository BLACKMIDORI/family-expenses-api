package com.blackmidori.apps.familyexpenses.api.factory;

import com.blackmidori.apps.familyexpenses.api.model.User;
import com.blackmidori.apps.familyexpenses.api.model.Workspace;
import lombok.NonNull;

import java.time.OffsetDateTime;

public class WorkspaceFactory {

    public Workspace create(@NonNull User user) {
        return new Workspace(user);
    }
}
