"use client";

import { useState } from "react";
import Link from "next/link";
import { MapPin, User, Mail, Lock, ArrowRight } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";

export default function SignupPage() {
    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError] = useState("");

    const { signup, isSigningUp } = useAuth();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        if (password !== confirmPassword) {
            setError("Passwords do not match");
            return;
        }

        signup({ name, email, password });
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-muted/30 p-4 py-8">
            <div className="w-full max-w-md space-y-8 animate-in fade-in zoom-in-95 duration-500">
                <div className="flex flex-col items-center justify-center text-center space-y-2 mb-4">
                    <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary shadow-sm mb-2">
                        <MapPin className="h-6 w-6" />
                    </div>
                    <h1 className="text-2xl font-bold tracking-tight">Join Circle</h1>
                </div>

                <Card className="border-muted shadow-lg shadow-primary/5">
                    <CardHeader className="space-y-1 pb-4">
                        <CardTitle className="text-xl font-bold">Create an account</CardTitle>
                        <CardDescription>
                            Enter your details below to create your account
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            {error && (
                                <div className="p-3 text-sm font-medium text-destructive bg-destructive/10 border border-destructive/20 rounded-md">
                                    {error}
                                </div>
                            )}

                            <div className="space-y-2">
                                <Label htmlFor="name" className="font-semibold text-sm">Full Name</Label>
                                <div className="relative">
                                    <User className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="name"
                                        placeholder="John Doe"
                                        className="pl-9 h-10 transition-all focus-visible:ring-primary/50"
                                        required
                                        value={name}
                                        onChange={(e) => setName(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="email" className="font-semibold text-sm">Email address</Label>
                                <div className="relative">
                                    <Mail className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="email"
                                        type="email"
                                        placeholder="name@example.com"
                                        className="pl-9 h-10 transition-all focus-visible:ring-primary/50"
                                        required
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="password" className="font-semibold text-sm">Password</Label>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="password"
                                        type="password"
                                        className="pl-9 h-10 transition-all focus-visible:ring-primary/50"
                                        required
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                    />
                                </div>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="confirmPassword" className="font-semibold text-sm">Confirm Password</Label>
                                <div className="relative">
                                    <Lock className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                                    <Input
                                        id="confirmPassword"
                                        type="password"
                                        className="pl-9 h-10 transition-all focus-visible:ring-primary/50"
                                        required
                                        value={confirmPassword}
                                        onChange={(e) => setConfirmPassword(e.target.value)}
                                    />
                                </div>
                            </div>

                            <Button
                                type="submit"
                                className="w-full h-11 mt-6 text-base font-semibold transition-all shadow-md hover:shadow-lg flex items-center justify-center gap-2 group"
                                disabled={isSigningUp}
                            >
                                {isSigningUp ? "Creating account..." : (
                                    <>
                                        Create Account
                                        <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-1" />
                                    </>
                                )}
                            </Button>
                        </form>
                    </CardContent>
                    <CardFooter className="flex flex-col items-center justify-center space-y-4 pt-4 border-t border-muted/50 bg-muted/10 rounded-b-xl px-6">
                        <div className="text-sm text-muted-foreground">
                            Already have an account?{" "}
                            <Link href="/login" className="text-primary font-semibold hover:underline">
                                Sign in instead
                            </Link>
                        </div>
                        <p className="text-xs text-center text-muted-foreground leading-relaxed">
                            By clicking continue, you agree to our{" "}
                            <Link href="#" className="underline underline-offset-4 hover:text-primary transition-colors">Terms of Service</Link>{" "}
                            and{" "}
                            <Link href="#" className="underline underline-offset-4 hover:text-primary transition-colors">Privacy Policy</Link>.
                        </p>
                    </CardFooter>
                </Card>
            </div>
        </div>
    );
}
