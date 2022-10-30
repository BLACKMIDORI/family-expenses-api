package com.blackmidori.apps.familyexpenses.api.model;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.data.mongodb.core.mapping.DBRef;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
public class PayerPaymentWeight {
    @DBRef(lazy = true)
    @NonNull
    private Payer payer;
    private double weight;
}
