package com.blackmidori.familyexpenses.api.application.converter;

import org.springframework.core.convert.converter.Converter;

import java.time.OffsetDateTime;

public class OffsetDateTimeWriteConverter implements Converter<OffsetDateTime, String> {
    @Override
    public String convert(OffsetDateTime offsetDateTime) {
        return offsetDateTime.toString();
    }

}