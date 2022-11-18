package com.blackmidori.familyexpenses.api.application.exception;

public class EntityNotFound extends RuntimeException {
    private final Class type;
    private final String entityId;

    public EntityNotFound(Class type, String entityId) {
        this.type = type;
        this.entityId = entityId;
    }

    public String getMessage() {
        return "incorrect id for " + type.getSimpleName() + ": " + entityId;
    }
}
