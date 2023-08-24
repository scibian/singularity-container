#
# Copyright (c) 2017-2022, SyLabs, Inc. All rights reserved.
# Copyright (c) 2017, SingularityWare, LLC. All rights reserved.
#
# Copyright (c) 2015-2017, Gregory M. Kurtzer. All rights reserved.
#
# Copyright (c) 2016, The Regents of the University of California, through
# Lawrence Berkeley National Laboratory (subject to receipt of any required
# approvals from the U.S. Dept. of Energy).  All rights reserved.
#
# This software is licensed under a customized 3-clause BSD license.  Please
# consult LICENSE file distributed with the sources of this project regarding
# your rights to use or distribute this software.
#
# NOTICE.  This Software was developed under funding from the U.S. Department of
# Energy and the U.S. Government consequently retains certain rights. As such,
# the U.S. Government has been granted for itself and others acting on its
# behalf a paid-up, nonexclusive, irrevocable, worldwide license in the Software
# to reproduce, distribute copies to the public, prepare derivative works, and
# perform publicly and display publicly, and to permit other to do so.
#
#

# Disable debugsource packages; otherwise it ends up with an empty %files
#   file in debugsourcefiles.list on Fedora
%undefine _debugsource_packages

Summary: Application and environment virtualization
Name: singularity-ce
Version: 3.10.5
Release: 1%{?dist}
# See LICENSE.md for first party code (BSD-3-Clause and LBNL BSD)
# See LICENSE_THIRD_PARTY.md for incorporated code (ASL 2.0)
# See LICENSE_DEPENDENCIES.md for dependencies
# License identifiers taken from: https://fedoraproject.org/wiki/Licensing
License: BSD-3-Clause and LBNL BSD and ASL 2.0
URL: https://www.sylabs.io/singularity/
Source: %{name}-3.10.5.tar.gz
ExclusiveOS: linux

BuildRequires: git
BuildRequires: gcc
BuildRequires: make
BuildRequires: libseccomp-devel
BuildRequires: glib2-devel
BuildRequires: cryptsetup

Requires: cryptsetup
Requires: glib2
Requires: runc
%if "%{_target_vendor}" == "suse"
Requires: squashfs
Requires: libseccomp2
%else
Requires: squashfs-tools
Requires: libseccomp
%endif

# there's no golang for ppc64, just ppc64le
ExcludeArch: ppc64

Provides: %{name}-runtime

# Conflicts with non-CE packages
Conflicts: singularity
# Conflicts with Apptainer, which installs the `/usr/bin/singularity` compatibility executable
Conflicts: apptainer
# Conflicts with SingularityPRO basic packaging (not other variants).
Conflicts: singularitypro24
Conflicts: singularitypro25
Conflicts: singularitypro26
Conflicts: singularitypro31
Conflicts: singularitypro35
Conflicts: singularitypro37
Conflicts: singularitypro39

%description
SingularityCE is the Community Edition of Singularity, an open source
container platform designed to be simple, fast, and secure.

%debug_package

%prep
%if "%{?buildroot}"
export RPM_BUILD_ROOT="%{buildroot}"
%endif

# Create our build root
chmod -R +w %{name}-%{version} || true
rm -rf %{name}-%{version}
mkdir %{name}-%{version}

%build
cd %{name}-%{version}

# Setup an empty GOPATH for the build
export GOPATH=$PWD/gopath
mkdir -p "$GOPATH"

# Extract the source
tar -xf "%SOURCE0"
cd %{name}-3.10.5

# Not all of these parameters currently have an effect, but they might be
#  used someday.  They are the same parameters as in the configure macro.
./mconfig -V %{version}-%{release} \
        --prefix=%{_prefix} \
        --exec-prefix=%{_exec_prefix} \
        --bindir=%{_bindir} \
        --sbindir=%{_sbindir} \
        --sysconfdir=%{_sysconfdir} \
        --datadir=%{_datadir} \
        --includedir=%{_includedir} \
        --libdir=%{_libdir} \
        --libexecdir=%{_libexecdir} \
        --localstatedir=%{_localstatedir} \
        --sharedstatedir=%{_sharedstatedir} \
        --mandir=%{_mandir} \
        --infodir=%{_infodir}

make -C builddir old_config=

%install
cd %{name}-%{version}

export GOPATH=$PWD/gopath
cd %{name}-3.10.5

make -C builddir DESTDIR=$RPM_BUILD_ROOT install

%files
%attr(4755, root, root) %{_libexecdir}/singularity/bin/starter-suid
%{_bindir}/singularity
%{_bindir}/run-singularity
%dir %{_libexecdir}/singularity
%{_libexecdir}/singularity/bin/starter
%{_libexecdir}/singularity/bin/conmon
%{_libexecdir}/singularity/cni/*
%dir %{_sysconfdir}/singularity
%config(noreplace) %{_sysconfdir}/singularity/*.conf
%config(noreplace) %{_sysconfdir}/singularity/*.toml
%config(noreplace) %{_sysconfdir}/singularity/*.json
%config(noreplace) %{_sysconfdir}/singularity/*.yaml
%config(noreplace) %{_sysconfdir}/singularity/global-pgp-public
%config(noreplace) %{_sysconfdir}/singularity/cgroups/*
%config(noreplace) %{_sysconfdir}/singularity/network/*
%config(noreplace) %{_sysconfdir}/singularity/seccomp-profiles/*
%dir %{_sysconfdir}/bash_completion.d
%{_sysconfdir}/bash_completion.d/*
%dir %{_localstatedir}/singularity
%dir %{_localstatedir}/singularity/mnt
%dir %{_localstatedir}/singularity/mnt/session
%{_mandir}/man1/singularity*
%license %{name}-%{version}/%{name}-3.10.5/LICENSE.md
%license %{name}-%{version}/%{name}-3.10.5/LICENSE_THIRD_PARTY.md
%license %{name}-%{version}/%{name}-3.10.5/LICENSE_DEPENDENCIES.md
%doc %{name}-%{version}/%{name}-3.10.5/README.md
%doc %{name}-%{version}/%{name}-3.10.5/CHANGELOG.md
%doc %{name}-%{version}/%{name}-3.10.5/CONTRIBUTING.md
%doc %{name}-%{version}/%{name}-3.10.5/CONTRIBUTORS.md

%changelog

