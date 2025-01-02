<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('store', function (Blueprint $table) {
            $table->id();
            $table->string("name");
            $table->text("description")->nullable();
            $table->string("location")->nullable();
            $table->decimal("latitude", 10, 8)->nullable();
            $table->decimal("longtitude", 10, 8)->nullable();
            $table->string("phone_number")->nullable();
            $table->string("email")->nullable();
            $table->string("image_url")->nullable();
            $table->string("whatsapp_link")->nullable();
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('store');
    }
};
