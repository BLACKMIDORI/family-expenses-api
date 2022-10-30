package com.blackmidori.apps.familyexpenses.api.dto;

import com.blackmidori.apps.familyexpenses.api.model.Charge;
import lombok.*;

import javax.validation.Valid;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;

@Data
public class ChargesRecordDto {
    @NotBlank
    private String chargesModelId;
    @NotEmpty
    private List<@Valid @NotNull ChargeDto> charges;
}
