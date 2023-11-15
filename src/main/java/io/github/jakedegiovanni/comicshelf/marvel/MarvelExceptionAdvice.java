package io.github.jakedegiovanni.comicshelf.marvel;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
@Slf4j
class MarvelExceptionAdvice {

    @ExceptionHandler({Exception.class})
    ResponseEntity<Void> anyException(Exception e) {
        log.error("got exception processing request", e);
        return ResponseEntity.internalServerError().build();
    }

    @ExceptionHandler({JsonProcessingException.class, JsonMappingException.class})
    ResponseEntity<Void> jsonException(Exception e) {
        log.error("could not handle json", e);
        return ResponseEntity.internalServerError().build();
    }
}
