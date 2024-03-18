<?php

namespace Comicshelf\Controllers;

use Phalcon\Mvc\Controller;

class IndexController extends Controller
{

    public function indexAction(): void
    {
        $this->view->name = "Phalcon App";
    }
}