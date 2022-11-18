package com.blackmidori.familyexpenses.api.dto;

import lombok.Data;

import javax.validation.Valid;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;

@Data
public class ChargeAssociationDto {
    @NotBlank
    private String expenseId;
    @NotEmpty
    private List<@Valid @NotNull PayerPaymentWeightDto> paymentWeights;
}
