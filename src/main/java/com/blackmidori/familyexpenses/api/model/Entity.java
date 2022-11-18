package com.blackmidori.familyexpenses.api.model;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.NonNull;
import org.springframework.data.annotation.Id;

import java.time.OffsetDateTime;

@Data
abstract public class Entity {
    @Id
    private String id;
    @NonNull
    private OffsetDateTime creationDateTime;

    Entity(){
        creationDateTime = OffsetDateTime.now();
    }
}
