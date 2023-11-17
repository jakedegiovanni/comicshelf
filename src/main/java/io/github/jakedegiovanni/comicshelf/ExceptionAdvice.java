package io.github.jakedegiovanni.comicshelf;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
@Slf4j
public class ExceptionAdvice {

    @ExceptionHandler({Exception.class})
    public ResponseEntity<Void> anyException(Exception e) {
        log.error("got exception processing request", e);
        return ResponseEntity.internalServerError().build();
    }

    @ExceptionHandler({JsonProcessingException.class, JsonMappingException.class})
    public ResponseEntity<Void> jsonException(Exception e) {
        log.error("could not handle json", e);
        return ResponseEntity.internalServerError().build();
    }
}
