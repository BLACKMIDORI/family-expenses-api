package com.blackmidori.familyexpenses.api.model;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
public class PayerPaymentAmount {
    @NonNull
    private Payer payer;
    private double amount;
}
