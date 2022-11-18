package com.blackmidori.familyexpenses.api.model;

import lombok.*;
import org.springframework.data.annotation.Id;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
@RequiredArgsConstructor
public class User {
    @Id
    private String id;
    @NonNull
    private String name;
    @NonNull
    private List<String> roles;

    static public User developmentUser = new User("DevelopmentUserInstance","Developer", List.of("admin"));
}
