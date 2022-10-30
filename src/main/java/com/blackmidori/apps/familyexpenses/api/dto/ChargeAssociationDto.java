package com.blackmidori.apps.familyexpenses.api.dto;

import com.blackmidori.apps.familyexpenses.api.model.PayerPaymentWeight;
import lombok.Data;
import lombok.NonNull;

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
