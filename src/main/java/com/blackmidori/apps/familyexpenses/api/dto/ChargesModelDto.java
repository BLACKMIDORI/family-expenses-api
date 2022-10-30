package com.blackmidori.apps.familyexpenses.api.dto;

import lombok.*;

import javax.validation.Valid;
import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;

@Data
public class ChargesModelDto {
    @NotEmpty
    private List<@Valid @NotNull ChargeAssociationDto> chargesAssociations;
    @NotBlank
    private String workspaceId;
}
