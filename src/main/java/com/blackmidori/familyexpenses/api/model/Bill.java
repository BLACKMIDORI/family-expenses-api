package com.blackmidori.familyexpenses.api.model;

import com.mongodb.lang.Nullable;
import lombok.NonNull;

import java.time.OffsetDateTime;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.RequiredArgsConstructor;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
public class Bill {
    @Nullable
    private OffsetDateTime dueDateTime;
    @NonNull
    private Expense expense;
    private double amount;
}
