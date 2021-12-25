<?php

namespace Tests\Browser\Debug;

use Livewire\Livewire;
use Tests\Browser\Init\Component;
use Tests\Browser\TestCase;

class Test extends TestCase
{
    public function test()
    {
        $this->browse(function ($browser) {
            Livewire::visit($browser, Component::class);

            sleep(1);
        });
    }
}
