package com.blackmidori.familyexpenses.api.model;

import lombok.*;
import org.springframework.data.mongodb.core.mapping.DBRef;

import java.time.OffsetDateTime;
import java.util.List;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class ChargesModel extends Entity {
    @NonNull
    private OffsetDateTime updatingDateTime;
    @NonNull
    private List<ChargeAssociation> chargesAssociations;
    //////// HAVE WORKSPACE ////////
    @DBRef(lazy = true)
    @NonNull
    private Workspace workspace;
}
