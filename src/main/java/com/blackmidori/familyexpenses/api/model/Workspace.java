package com.blackmidori.familyexpenses.api.model;

import lombok.*;
import org.springframework.data.mongodb.core.mapping.DBRef;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class Workspace extends Entity {
    @DBRef(lazy = true)
    @NonNull
    private User owner;
}
