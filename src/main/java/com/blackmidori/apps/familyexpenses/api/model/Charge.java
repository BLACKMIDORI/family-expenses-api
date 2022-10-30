package com.blackmidori.apps.familyexpenses.api.model;

import lombok.NonNull;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.RequiredArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
public class Charge {
    @NonNull
    private Bill bill;

    @NonNull
    private List<PayerPaymentAmount> paymentAmountList;

}
