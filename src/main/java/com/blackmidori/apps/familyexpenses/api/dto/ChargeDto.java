package com.blackmidori.apps.familyexpenses.api.dto;

import lombok.Data;

import javax.validation.Valid;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;

@Data
public class ChargeDto {
    @Valid
    @NotNull
    private BillDto bill;
    @NotEmpty
    private List<@Valid @NotNull PayerPaymentAmountDto> paymentAmountList;

}
