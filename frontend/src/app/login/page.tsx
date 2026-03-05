"use client";

import { useState } from "react";
import Link from "next/link";
import { MapPin, Mail, Lock, ArrowRight } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";

export default function LoginPage() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const { login, isLoggingIn } = useAuth();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        login({ email, password });
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-muted/30 p-4">
            <div className="w-full max-w-md space-y-8 animate-in fade-in zoom-in-95 duration-500">
                <div className="flex flex-col items-center justify-center text-center space-y-2 mb-8">
                    <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-primary/10 text-primary shadow-sm mb-2">
                        <MapPin className="h-8 w-8" />
                    </div>
                    <h1 className="text-3xl font-bold tracking-tight">Welcome to Circle</h1>
                    <p className="text-muted-foreground font-medium">Connect with people around you.</p>
                </div>

                <Card className="border-muted shadow-lg shadow-primary/5">
                    <CardHeader className="space-y-1 pb-6">
                        <CardTitle className="text-2xl font-bold">Login</CardTitle>
                        <CardDescription className="text-base">
                            Enter your credentials to access your local feed
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="email" className="font-semibold">Email address</Label>
                                <div className="relative">
                                    <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="email"
                                        type="email"
                                        placeholder="name@example.com"
                                        className="pl-9 h-11 transition-all focus-visible:ring-primary/50"
                                        required
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                    />
                                </div>
                            </div>
                            <div className="space-y-2">
                                <div className="flex items-center justify-between">
                                    <Label htmlFor="password" className="font-semibold">Password</Label>
                                    <Link href="#" className="text-sm font-medium text-primary hover:underline">
                                        Forgot password?
                                    </Link>
                                </div>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="password"
                                        type="password"
                                        className="pl-9 h-11 transition-all focus-visible:ring-primary/50"
                                        required
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                    />
                                </div>
                            </div>

                            <Button
                                type="submit"
                                className="w-full h-11 mt-6 text-base font-semibold transition-all shadow-md hover:shadow-lg flex items-center justify-center gap-2 group"
                                disabled={isLoggingIn}
                            >
                                {isLoggingIn ? "Signing in..." : (
                                    <>
                                        Sign In
                                        <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
                                    </>
                                )}
                            </Button>
                        </form>
                    </CardContent>
                    <CardFooter className="flex flex-col items-center justify-center space-y-4 pt-4 border-t border-muted/50 bg-muted/10 rounded-b-xl">
                        <div className="text-sm text-muted-foreground">
                            Don't have an account?{" "}
                            <Link href="/signup" className="text-primary font-semibold hover:underline">
                                Create new account
                            </Link>
                        </div>
                    </CardFooter>
                </Card>
            </div>
        </div>
    );
}
