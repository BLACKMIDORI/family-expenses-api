package com.blackmidori.familyexpenses.api.dto;

import lombok.Data;

import javax.validation.constraints.NotBlank;

@Data
public class ExpenseDto  {
    @NotBlank
    private String name;
    @NotBlank
    private String workspaceId;
}