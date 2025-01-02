<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Store extends Model
{
    use HasFactory;

    protected $table = 'store';

    protected $fillable = [
        'name',
        'description',
        'location',
        'latitude',
        'longtitude',
        'phone_number',
        'email',
        'image_url',
        'whatsapp_link'
    ];
}
