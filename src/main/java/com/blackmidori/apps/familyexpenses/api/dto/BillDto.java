package com.blackmidori.apps.familyexpenses.api.dto;

import lombok.Data;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import java.time.OffsetDateTime;

@Data
public class BillDto {
    @NotNull
    private OffsetDateTime dueDateTime;
    @NotBlank
    private String expenseId;
    @NotNull
    private Double amount;
}
