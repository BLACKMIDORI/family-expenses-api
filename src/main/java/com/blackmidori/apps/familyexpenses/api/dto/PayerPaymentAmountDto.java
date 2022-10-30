package com.blackmidori.apps.familyexpenses.api.dto;

import lombok.Data;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;

@Data
public class PayerPaymentAmountDto {
    @NotBlank
    private String payerId;
    @NotNull
    private Double amount;
}
