package com.blackmidori.apps.familyexpenses.api.model;

import lombok.*;
import org.springframework.data.mongodb.core.mapping.DBRef;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class Payer extends Entity{
    @NonNull
    private String name;
    //////// HAVE WORKSPACE ////////
    @DBRef(lazy = true)
    @NonNull
    private Workspace workspace;
}
