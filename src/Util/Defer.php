<?php

namespace Comicshelf\Util;

class Defer
{
    private $fn;

    public function __construct($fn)
    {
        $this->fn = $fn;
    }

    public function __destruct()
    {
        call_user_func($this->fn);
    }
}
