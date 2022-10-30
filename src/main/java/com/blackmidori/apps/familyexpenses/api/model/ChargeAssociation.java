package com.blackmidori.apps.familyexpenses.api.model;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.data.mongodb.core.mapping.DBRef;

import java.util.List;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
public class ChargeAssociation {
    @DBRef(lazy = true)
    @NonNull
    private Expense expense;
    @NonNull
    private List<PayerPaymentWeight> paymentWeights;
}
