package com.blackmidori.familyexpenses.api.model;

import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.DBRef;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
public class AuthenticationMethod {
    @Id
    private String id;
    @DBRef(lazy = true)
    @NonNull
    private User user;
}
