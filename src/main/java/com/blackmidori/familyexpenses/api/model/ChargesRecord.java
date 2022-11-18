package com.blackmidori.familyexpenses.api.model;

import lombok.*;

import java.time.OffsetDateTime;
import java.util.List;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
@EqualsAndHashCode(callSuper = true)
public class ChargesRecord extends Entity {
    @NonNull
    private OffsetDateTime openingDateTime;
    private OffsetDateTime closingDateTime;
    @NonNull
    private ChargesModel chargesModel;

    @NonNull
    private List<Charge> charges;
}
