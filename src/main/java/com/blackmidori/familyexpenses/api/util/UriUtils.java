package com.blackmidori.familyexpenses.api.util;

import org.springframework.web.servlet.support.ServletUriComponentsBuilder;

import java.net.URI;

public class UriUtils {
    public static URI getCreatedUrl(String resourceId){
        return ServletUriComponentsBuilder.fromCurrentRequest().pathSegment(resourceId).build().toUri();
    }
}
